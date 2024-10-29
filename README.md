# asdf

## example

**part1_bass96**
```
.steps96 c.vel50.gate50
a4
f3
g
```


**part2_bass**
```
(c c) * 4 d
a - - a
f
g
```

**part2_chords**
```
Cmaj7.arp.up4.down4
Amin/E;4
F
G
```


**clock**
```
(c8 c0) * 8
```

### play

*crow.voct4*
```
clock
```

*crow.voct1.env2.bpm120*
```
part1_chords * 4
part2_chords * 4
```

*crow.voct3*
```
part1_bass * 4
part2_bass * 4
```