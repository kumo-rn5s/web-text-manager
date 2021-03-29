# ECS-Exec の使用方法

1. SSM-Agent

	[SSM-Agentについて](https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/ssm-agent.html)

2. aws-cliの更新

	```bash
	aws --version
	aws-cli/2.1.31 Python/3.8.8 Darwin/19.6.0 exe/x86_64 prompt/off
	```
	[MacOS上aws-cliのアップデート](https://docs.aws.amazon.com/ja_jp/cli/latest/userguide/install-cliv2-mac.html)
	
	```bash
	curl "https://awscli.amazonaws.com/AWSCLIV2-2.0.30.pkg" -o "AWSCLIV2.pkg"
	sudo installer -pkg AWSCLIV2.pkg -target /
	```

3. Session Manager Install

	[Download Page](https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html)

	```bash
	curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/mac/sessionmanager-bundle.zip" -o "sessionmanager-bundle.zip"
	unzip sessionmanager-bundle.zip
	sudo ./sessionmanager-bundle/install -i /usr/local/sessionmanagerplugin -b /usr/local/bin/session-manager-plugin
	
	session-manager-plugin
	// The Session Manager plugin was installed successfully. Use the AWS CLI to start a session.
	
	// そのあと、logレベルをinfo/debugにする
	https://docs.aws.amazon.com/ja_jp/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html#install-plugin-verify
	```

4. ECS Execの実装

	1.  SSMセッションマネージャー関連の権限を付与したIAMロール (ECSタスクロール) を用意する
	
		```yaml
		- Effect: Allow
	            Action:
	              - ssmmessages:CreateControlChannel
	              - ssmmessages:CreateDataChannel
	              - ssmmessages:OpenControlChannel
	              - ssmmessages:OpenDataChannel
	            Resource: "*"
		```

	2. ECSサービスで「enableExecuteCommand」の設定を有効にする
	
		```yaml
		Type: AWS::ECS::Service
		      EnableExecuteCommand: true
		```
	
	3. (ECS Execを使ってコンテナ上でコマンドを実行する)
	
		```bash
		aws ecs execute-command \
	    --cluster <cluster name> \
	    --task <task ID> \
	    --interactive \
	    --command "<command lines>"
	    
	    // Examples
	    
	    aws ecs execute-command \
	    	--cluster cfn-test-ecs-cluster \
	    	--task 8e55260c92a942408efad76509e10604 \
	    	--interactive \
	    	--command "ls /var/www/html"
		```
		
		
	
