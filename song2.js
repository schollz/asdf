bpm = 160
part1 = `
.bpm${bpm}.gate50
c4 a d5 
c5 e g5
d5 g5 c5
c6 e f4 c6 e f4
c2 b1 g2 g1
`

part2 = `.gate50.bpm120.velocity30.reverb0
C;3.arp.up8.down6.up2.down4
Em;3.arp.up8.down6.up2.down4
Am/E;3.arp.up8.down6.up2.down4
F;3.arp.up8.down6.up2.down4
`

part3 = `c5.velocity120.gate50
d e
c5
g f
`

supercollider_polyperc_1 = "part2"

supercollider_polyperc_2 = "part3"

