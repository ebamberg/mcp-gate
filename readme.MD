# MCP Gate
MCP Gateway and proxy

⚠️ **This project is NOT intended for production use.** The project is Proof-of-concept phase ⚠️

# Getting started

start the mcp-gate server

```
mcp-gate server
```

## Install in Claude Desktop

!! only MacOS and Windows 

you can install mcp-gate into an installed local Claude Desktop using the tools cli.

```
mcp-gate install claude
```

This will read Claud Desktop config file and adds mcp-gate as a mcp-server to Claude Desktop. 
After restarting Claude Desktop mcp-gate is visible as a tool. 

Running the command again will reinstall and override previous mcp-gate config. This is useful if the folder to the exe has changed.

Before writing the Claude Desktop config file a backup is created in the same folder as the original config file lives.

# Command line summary

| command | description                                                              |
|---------|--------------------------------------------------------------------------|
| server  | start the gateway & proxy in mcp server mode                             |
| install | installs the gateway in target for example `install claude`              |

# Disclaimer
This project is provided as-is and it's still in experimental phase. The authors make no warranties regarding the code's suitability for any particular purpose. Users should not deploy this code in production environments without significant modifications and testing.