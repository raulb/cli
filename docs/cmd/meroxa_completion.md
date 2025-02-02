## meroxa completion

Generate completion script

### Synopsis

To load completions:

	Bash:
	
	$ source <(meroxa completion bash)
	
	# To load completions for each session, execute once:
	Linux:
	  $ meroxa completion bash > /etc/bash_completion.d/meroxa
	MacOS:
	  $ meroxa completion bash > /usr/local/etc/bash_completion.d/meroxa
	
	Zsh:
	
	# If shell completion is not already enabled in your environment you will need
	# to enable it.  You can execute the following once:
	
	$ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
	# To load completions for each session, execute once:
	$ meroxa completion zsh > "${fpath[1]}/_meroxa"
	
	# You will need to start a new shell for this setup to take effect.
	
	Fish:
	
	$ meroxa completion fish | source
	
	# To load completions for each session, execute once:
	$ meroxa completion fish > ~/.config/fish/completions/meroxa.fish
	

```
meroxa completion [bash|zsh|fish|powershell]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/meroxa.env)
      --debug           display any debugging information
      --json            output json
```

### SEE ALSO

* [meroxa](meroxa.md)	 - The Meroxa CLI

