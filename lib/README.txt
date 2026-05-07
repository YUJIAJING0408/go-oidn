Place dynamic libraries for different platforms in the lib folder.
For example, on the Windows platform, download the Windows executable file library of Oidn and copy all the DLL files in the bin directory to 'lib\windows\'
Either manually set the global environment variable OIDN_LIB_PATH or set it at runtime, with the priority order being: oidn.SetLibraryPath > OIDN_LIB_PATH > lib\*\*