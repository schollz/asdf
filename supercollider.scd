(
s.waitForBoot({
	var server=Server.default;
	~notes = Dictionary.new();
	~params = Dictionary.new();

	SynthDef("polyperc", {
		arg outReverb=4, outDry=5, note=60, velocity=120, gate=1, pan=0,
		attack=0.05,decay=0.1,sustain=0.5,release=1,lpf=1000, reverb = 20;
		var snd = Pulse.ar(note.midicps, 0.3);
		snd = MoogFF.ar(snd,Lag.kr(lpf,2),2);
		snd = snd * EnvGen.ar(Env.adsr(attack,decay,sustain,release),gate,doneAction:2);
		snd = Pan2.ar(snd,pan/100.0);
		snd = snd * LinLin.kr(Lag.kr(velocity,2),0,127,0,2);
		//FreeSelf.kr(TDelay.kr(DC.kr(1),10));
		Out.ar(outReverb,snd*(reverb/100));
		Out.ar(outDry, snd*(1-(reverb/100)));
	}).add;
	SynthDef("main",{
		var out=0, inReverb=4, inDry=6;
		var sndDry = In.ar(inDry,2);
		var sndReverb = In.ar(inReverb,2);
		var snd;
		snd = sndDry + Fverb.ar(sndReverb[0],sndReverb[1],50,
			tail_density: LFNoise2.kr(1/3).range(50,90),
			decay: LFNoise2.kr(1/3).range(50,90),
			input_diffusion_1: LFNoise2.kr(1/3).range(50,90),
		);
		Out.ar(out, snd*4);
	}).add;

	if (~busReverb.isNil,{
		~busReverb = Bus.audio(server,2);
	});
	if (~busDry.isNil,{
		~busDry = Bus.audio(server,2);
	});
	server.sync;

	if (~main.isNil,{
		~main = Synth.tail(server,"main",[
			\inReverb,~busReverb.postln,
			\inDry,~busDry.postln,
		]);
		server.sync;
	});


	~fnSet = {
		arg id, synthdef, key, value;
		var keyID=id.asString++synthdef;
		if (~params.at(keyID).isNil,{
			~params.put(keyID,Dictionary.new());
		});
		~params.at(keyID).put(key,value);
		// loop through all synths
		if (~notes.at(keyID).notNil,{
			~notes.at(keyID).keysValuesDo({ arg k, v;
				if (v.isRunning,{
					v.set(key.asSymbol,value);
				});
			});

		});

	};
	~fnNoteOn = {
		arg id,synthdef,note,velocity;
		var keyID=id.asString++synthdef;
		var args=[
			\note,note,
			\velocity,velocity,
			\outReverb, ~busReverb,
			\outDry, ~busyDry,
		];
		if (~params.at(keyID).notNil,{
			~params.at(keyID).keysValuesDo({ arg k, v;
				args=args++[k,v];
			});

		});
		if (~notes.at(keyID).isNil,{
			~notes.put(keyID,Dictionary.new());
		});
		if (~notes.at(keyID).at(note).notNil,{
			if (~notes.at(keyID).at(note).isRunning,{
				~notes.at(keyID).at(note).set(\gate,0);
			});
		});
		[keyID,note,"playing"].postln;
		~notes.at(keyID).put(note,Synth.head(server,synthdef,args).onFree({
			[keyID,note,"finished"].postln;
		}));
		NodeWatcher.register(~notes.at(keyID).at(note));
	};
	~fnNoteOff = {
		arg id, synthdef, note;
		var keyID=id.asString++synthdef;
		if (~notes.at(keyID).isNil,{
			~notes.put(keyID,Dictionary.new());
		});
		if (~notes.at(keyID).at(note).notNil,{
			if (~notes.at(keyID).at(note).isRunning,{
				~notes.at(keyID).at(note).set(\gate,0);
			});
		});

	};

	if (~osc.isNil,{
		"setting up osc".postln;
		~osc = OSCFunc({ arg msg, time, addr, recvPort;
			var msgtype = msg[1].asString;
			var id = msg[2].asInteger;
			var synthdef = msg[3].asString;
			if (msgtype=="note_on",{
				var note=msg[4].asInteger;
				var velocity=msg[5].asInteger;
				// [synthdef++id.asString,"note_on",note,velocity].postln;
				~fnNoteOn.(id,synthdef,note,velocity);
			});
			if (msgtype=="note_off",{
				var note=msg[4].asInteger;
				// [synthdef++id.asString,"note_off",note].postln;
				~fnNoteOff.(id,synthdef,note);
			});
			if (msgtype=="set",{
				var paramName=msg[4].asString;
				var paramValue=msg[5].asFloat;

				[synthdef++id.asString,"set",paramName,paramValue].postln;
				~fnSet.(id,synthdef,paramName,paramValue);
			});

		}, '/asdf', recvPort:7771); // create
	});
});
)

NetAddr("127.0.0.1",7771).sendMsg("/asdf", "there");

~fnSet.(1,"polyperc","lpf",4000);
~fnSet.(1,"polyperc","attack",0.2);
~fnNoteOn.(1,"polyperc",72,120);
~fnNoteOff.(1,"polyperc",72);