(
s.waitForBoot({
	var server=Server.default;
	~notes = Dictionary.new();
	~params = Dictionary.new();
	~messageQueue = List.new;
	SynthDef("polyperc", {
		arg outReverb=4, outDry=5, note=60, velocity=120, gate=1, pan=0,
		attack=0.05,decay=0.1,sustain=0.5,release=1,lpf=1000, reverb = 0;
		var snd = Pulse.ar(note.midicps, 0.3);
		snd = MoogFF.ar(snd,Lag.kr(lpf,2),2);
		snd = snd * EnvGen.ar(Env.adsr(attack,decay,sustain,release),gate-EnvGen.kr(Env.new([0,0,1],[3,0.1]),1),doneAction:2);
		snd = Pan2.ar(snd,pan/100.0);
		snd = snd * LinLin.kr(Lag.kr(velocity,2),0,127,0,2);
		Out.ar(outReverb,snd*(reverb/100));
		Out.ar(outDry, snd*(1-(reverb/100)));
	}).add;

	SynthDef("jp",{
		arg outReverb=4, outDry=5, note=60, velocity=120, gate=1, pan=0,
		attack=0.05,decay=0.1,sustain=0.5,release=1,lpf=10, reverb = 0,
		amp=0.75,db=0,mix=100, detune = 40;
		var snd;
		var freq=note.midicps;
		var detuneCurve = { |x|
			(10028.7312891634*x.pow(11)) -
			(50818.8652045924*x.pow(10)) +
			(111363.4808729368*x.pow(9)) -
			(138150.6761080548*x.pow(8)) +
			(106649.6679158292*x.pow(7)) -
			(53046.9642751875*x.pow(6)) +
			(17019.9518580080*x.pow(5)) -
			(3425.0836591318*x.pow(4)) +
			(404.2703938388*x.pow(3)) -
			(24.1878824391*x.pow(2)) +
			(0.6717417634*x) +
			0.0030115596
		};
		var centerGain = { |x| (-0.55366 * x) + 0.99785 };
		var sideGain = { |x| (-0.73764 * x.pow(2)) + (1.2841 * x) + 0.044372 };

		var center = Mix.new(SawDPW.ar(freq, Rand()));
		var detuneFactor = freq * detuneCurve.(detune/100.0);
		var freqs = [
			(freq - (detuneFactor * 0.11002313)),
			(freq - (detuneFactor * 0.06288439)),
			(freq - (detuneFactor * 0.01952356)),
			// (freq + (detuneFactor * 0)),
			(freq + (detuneFactor * 0.01991221)),
			(freq + (detuneFactor * 0.06216538)),
			(freq + (detuneFactor * 0.10745242))
		];
		var side = Mix.fill(6, { |n|
			SawDPW.ar(freqs[n], Rand(0, 2))
		});


		var sig =  (center * centerGain.(mix/100.0)) + (side * sideGain.(mix/100.0));
		sig = HPF.ar(sig ! 2, freq);
		sig = BLowPass.ar(sig,freq*LFNoise2.kr(1).range(4,20),1/0.707);
		sig = Pan2.ar(sig);
		snd = sig * amp * 12.neg.dbamp*db.dbamp;
		snd = snd * EnvGen.ar(Env.adsr(attack,decay,sustain,release),gate-EnvGen.kr(Env.new([0,0,1],[10,0.1]),1),doneAction:2);
		snd = Pan2.ar(snd,pan/100.0);
		snd = snd * LinLin.kr(Lag.kr(velocity,2),0,127,0,2);
		Out.ar(outReverb,snd*(reverb/100));
		Out.ar(outDry, snd*(1-(reverb/100)));
	}).add;

	Buffer.read(s,thisProcess.nowExecutingPath.dirname++"/data/909_1.wav",action:{ arg buf;
		buf.postln;
		SynthDef("volcano",{
			arg outReverb=4, outDry=5, note=60, velocity=120, gate=1, pan=0,
			attack=0.005,decay=0.1,sustain=0.5,release=2,lpf=1000, reverb = 0;
			var snd = PlayBuf.ar(1,buf);
			var snd1 = HPF.ar(snd,90);
			var snd1r, snd1l, snd1rl;
			// snd1 = (snd1*20).tanh;
			snd1 = AnalogVintageDistortion.ar(snd1,0.2,0.0,0.1,0.2);

			// reverb predelay time :
			snd1r = DelayN.ar(snd1, 0.02 + 0.03.rand);
			// 7 length modulated comb delays in parallel :
			snd1r = Mix.ar(Array.fill(7,{ CombL.ar(snd1r, 0.1, LFNoise1.kr(0.001.rand, 0.04, 0.05), 15.rand) }));
			// two parallel chains of 4 allpass delays (8 total) :
			4.do({ snd1r = AllpassN.ar(snd1r, 0.050, [0.050.rand, 0.050.rand], 1) });

			// reverb predelay time :
			snd1l = DelayN.ar(snd1, 0.02 + 0.03.rand);
			// 7 length modulated comb delays in parallel :
			snd1l = Mix.ar(Array.fill(7,{ CombL.ar(snd1l, 0.1, LFNoise1.kr(0.001.rand, 0.04, 0.05), 15.rand) }));
			// two parallel chains of 4 allpass delays (8 total) :
			4.do({ snd1l = AllpassN.ar(snd1r, 0.050, [0.050.rand, 0.050.rand], 1) });

			snd1 = Pan2.ar(snd1);

			snd1 = SelectX.ar(0.3,[snd1,[snd1l,snd1r]]);
			snd1 = snd1 * (1 - Amplitude.ar(snd));

			snd = Pan2.ar(snd,0);

			snd = (snd/2) + (snd1/2);

			snd = snd.tanh;

			snd = snd + CombC.ar(snd);

			snd = snd * EnvGen.ar(Env.new([1,1,0],[3,3],\sine));

			DetectSilence.ar(snd,doneAction:2);
			snd = snd * EnvGen.ar(Env.adsr(attack,decay,sustain,release),gate-EnvGen.kr(Env.new([0,0,1],[3,0.1]),1),doneAction:2);
			snd = Pan2.ar(snd,pan/100.0);
			snd = snd * LinLin.kr(Lag.kr(velocity,2),0,127,0,2);
			Out.ar(outReverb,snd*(reverb/100));
			Out.ar(outDry, snd*(1-(reverb/100)));

		}).add;

	});

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
		snd = AnalogTape.ar(snd,0.9,0.9,0.6,2);
		Out.ar(out, snd);
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
			~notes.at(keyID).keysValuesDo({ arg k1, d1;
				d1.keysValuesDo({ arg k,v;
					if (v.isRunning,{
						v.set(key.asSymbol,value);
					});
				});
			});

		});

	};
	~fnNoteOn = {
		arg id,synthdef,note,velocity;
		var keyID=id.asString++synthdef;
		var synth;
		var noteID;
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
		if (~notes.at(keyID).at(note).isNil,{
			~notes.at(keyID).put(note,Dictionary.new());
		});

		~notes.at(keyID).at(note).keysValuesDo({ arg k, v;
			if (v.isRunning,{
				v.set(\gate,0);
			});
		});

		noteID =100000000.rand;
		synth = Synth.head(server,synthdef,args).onFree({
			[keyID,note,"finished"].postln;
			~notes.at(keyID).at(note).removeAt(noteID);
		});
		~notes.at(keyID).at(note).put(noteID,synth);
		NodeWatcher.register(synth);
	};
	~fnNoteOff = {
		arg id, synthdef, note;
		var keyID=id.asString++synthdef;
		if (~notes.at(keyID).isNil,{
			~notes.put(keyID,Dictionary.new());
		});
		if (~notes.at(keyID).at(note).isNil,{
			~notes.at(keyID).put(note,Dictionary.new());
		});

		~notes.at(keyID).at(note).keysValuesDo({ arg k, v;
			if (v.isRunning,{
				v.set(\gate,0);
			});
		});

	};


	AppClock.sched(0.01, {
		Routine{
			// copy messageQueue
			var msgQ = List.new(~messageQueue.size);
			~messageQueue.do({ |v|
				msgQ.add(v)
			});
			~messageQueue.clear;
			msgQ.do { |msg|
				var msgtype = msg[1].asString;
				var id = msg[2].asInteger;
				var synthdef = msg[3];
				server.sync;
				if (msgtype=="note_on",{
					var note=msg[4].asInteger;
					var velocity=msg[5].asInteger;
					[synthdef++id.asString,"note_on",note,velocity].postln;
					~fnNoteOn.(id,synthdef,note,velocity);
				});
				if (msgtype=="note_off",{
					var note=msg[4].asInteger;
					[synthdef++id.asString,"note_off",note].postln;
					~fnNoteOff.(id,synthdef,note);
				});
				if (msgtype=="set",{
					var paramName=msg[4].asString;
					var paramValue=msg[5].asFloat;

					[synthdef++id.asString,"set",paramName,paramValue].postln;
					~fnSet.(id,synthdef,paramName,paramValue);
				});

			};
		}.play;
		0.01; // Re-schedule in 10 ms
	});

	if (~osc.isNil,{
		"setting up osc".postln;
		~osc = OSCFunc({ arg msg, time, addr, recvPort;
			~messageQueue.add(msg);
		}, '/asdf', recvPort:7771); // create
	});
});
)

// testing
(
~fnNoteOn.(0,"volcano",60,60);
)
(
~fnNoteOff.(0,"volcano",60,60);

)
(
~fnNoteOn.(0,"jp",30,60);
)
(
~fnNoteOff.(0,"jp",30,60);
)


