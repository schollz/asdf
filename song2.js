part3 = `
c4.bpm40.attack70.release100.velocity50.gate50.reverb70.detune70 d
e4.detune30,40,80 d.mix30,70
d4.detune50 g4.mix20,90.velocity60
c5 - b4
`

chords = `
Am;2.velocity40.gate70,60,30,50.release200,150,170.reverb40.attack50
F/A;2.gate70,60,30,50.release200,150,170
C;2.gate30,50,60,30,50.release140,110,170
Em/B;2.gate80,90,60,30,50.release70,200,150,170
`

kick = `
.bpm30.reverb30
c4.prob90.velocity40,35,50.release50,100,200.gate50,100 -  ~
`

supercollider_volcano="kick"
supercollider_jp_1 = "chords"
supercollider_jp_2 = "part3"
