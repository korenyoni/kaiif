Date: 2017-09-10
Title: Weird compiler errors
cat: ops,dev

The other day I ran into a compiler error which did not come with any explanation:

`c++: internal compiler error: Bus error (program cc1plus)`

If it had anything to do with build dependencies or a problem in the code (unlikely, since it's part of a heavily reviewed repository), g++ would spew something out.

It's likely you're running out of memory.

Fortunately my Cloud Provider (Joyent) allows you to resize vertically without even having to reboot your machine.

I gradually resized my instance and tried compiling again. Eventually I had enough RAM that the issue was solved.
