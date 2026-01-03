# DMX-Lighting Garbage (DLG)

I am developing DLG to learn about DMX. It is written in GoLang — a language that I am still new to as well.

DLG currently supports OpenDMX (via serial).
DLG offers raw DMX controls (channel control) and device-level controls.

Devices can be created based on [DeviceDefinitionFiles](https://wiki-de.dmxcontrol-projects.org/index.php?title=DDF_DMXC3) (DDF) in XML format.
Note that this format is not perfectly supported, so only basic functionality will work out of the box (for now).

This project is heavily inspired by the free software [DMXControl](https://www.dmxcontrol.de/).
