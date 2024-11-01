bpm = 160
part1 = `
.bpm${bpm}.gate50
c4 a d5 
c5 e g5
d5 g5 c5
c6 e f4 c6 e f4
c2 b1 g2 g1
`

part2 = `.gate50.bpm180
C;3.arp.up4.down4.attack1.reverb0
`
supercollider_polyperc = "part2"
