part1 = `
.bpm160.gate50.beats2.reverb0
c4.lpf1000.tr2,0,2,-2,12 d.tr2,-4 e.tr0,0,1 f
c.tr4,12,-12,2,0,2,-2,12 d.tr2,-4,8 e.tr0,0,1 f.tr5,-1,12,-12
`

part2 = `
.bpm120.gate90
Cmaj/B;3.arp.u12,4,8.d8
Em/B;3.arp.u6,9,12.d8.u12.d4
Am/E;3.arp.u4.d6.u9.d5
F/C;3.arp.u6.d2.t1.u8.d4
`
bpm=150
octave=4
part3 = `.bpm${bpm}
Cmaj/B;${octave}.arp.up8.down8
Em/B;${octave}.arp.up8.down8
Am/E;${octave}.arp.up8.down8
F/C;${octave}.arp.up8.down8
`
part4 = `.bpm${bpm}
Cmaj/B;${octave}.arp.up8.down8
`

clock=`.bpm${bpm}.beat1
(c3 c7) * 8`

//part2 = "a.bpm90.attack10 b c d"
// midi_op1 = "part1"

crow_voct1 = "clock"
crow_voct2 = "part3"
//supercollider_polyperc_1="part1"
//debug = "part3"
