Date: 2017-10-24
Title: Swapping the Chromebook search button
cat: ops

A very common configuration for programmers, especially among those who use Vim-keybindings, is the swappin of CapsLock and Escape:

![Swap caps lock and Escape](/images/keyboard-regular-swapped.png)

A modern way to do this on a Linux machine running Xorg is by using `setxkbmap`:

```
$ setxkbmap -option caps:swapescape
```

This is useless for a Chromebook, which requires a keyboard configuration like this:

![Swap Chromebook super and Escape](/images/keyboard-cb-swapped.png)

We know from using `xev` that the button in question has a keysym (the code used by xkb) of Super_L.

Posts such as [this](https://unix.stackexchange.com/a/65600) recommend using xkbcomp to modify your loaded xkb keyboard. However xkbcomp is becoming a bit [dated like xmodmap](https://wiki.archlinux.org/index.php/X_KeyBoard_extension#Using_keymap_.28deprecated.29), but more importantly [it can fail in startup scripts](https://askubuntu.com/q/437584), for me included --for whatever reason.

If we show the current xkb keyboard being loaded:

```
$ setxkbmap -print

xkb_keymap {
	xkb_keycodes  { include "evdev+chromebook_m(media)+aliases(qwerty)"	};
	xkb_types     { include "complete"	};
	xkb_compat    { include "complete+chromebook"	};
	xkb_symbols   { include "pc+us+inet(evdev)+chromebook_m_ralt(overlay)"	};
	xkb_geometry  { include "pc(pc104)"	};
};
```


We will see our loaded symbols. Let's grep the most specific symbol:

```
$ grep -rnw -e 'chromebook_m_ralt(overlay)'

rules/evdev:872:  chromebook_m_ralt        =   +inet(evdev)+chromebook_m_ralt(overlay)
rules/evdev:876:  chromebook_m_falco_ralt  =   +inet(evdev)+chromebook_m_ralt(overlay)
```

We're going to add a custom symbol to the first location in `evdev`, because it is most similar to our loaded symbols `inet(evdev)+chromebook_m_ralt(overlay)`:

1. Create a symbol file `symbols/super`:

```
// Swap Super and Escape
partial modifier_keys
xkb_symbols "swap_escape" {
        replace key <ESC> { [ Super_L ] };
        replace key <LWIN> { [ Escape ] };
};
```

2. Add the symbol to the `option = symbols` section:

```
super:swap_escape    =    +super(swap_escape)
```

3. Add the symbol to the line that was grepped:

```
chromebook_m_ralt   =   +inet(evdev)+chromebook_m_ralt(overlay)+super(swap_escape)
```

4. rebuild xkb-data:

```
dpkg-reconfigure xkb-data
```

Log out and log in. You should now have the new symbol loaded:

```
$ setxkbmap -print

xkb_keymap {
	xkb_keycodes  { include "evdev+chromebook_m(media)+aliases(qwerty)"	};
	xkb_types     { include "complete"	};
	xkb_compat    { include "complete+chromebook"	};
	xkb_symbols   { include "pc+us+inet(evdev)+chromebook_m_ralt(overlay)+super(swap_escape)"	};
	xkb_geometry  { include "pc(pc104)"	};
};
```

Your vim-friendly configuration is now working.
