# Map Key

An example of a 3x3 room map looks like below:
```
 |  |  | 
- -- -- -
 |  |  | 
 | ^|  | 
- --#-- -
 | v|  | 
 |  |  | 
- -- -- -
 |  |  | 
```

Each room takes up a 3x3 character square. All 6 exits are printed for each room
independantly of other rooms. The top left character is `^` if there is an up
exit. The bottom left character is `v` if there is a down exit. The `#`
indicates the player.