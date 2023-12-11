
```bash
docker build -t workday-local-redis -f Dockerfile-redis .

```

```bash
# get ec2 instance
aws ec2 describe-instances --filters "Name=instance-state-name,Values=running" --query 'Reservations[*].Instances[*].InstanceId' --output text

# get Public IP by EC2 instance
# aws ec2 describe-instances --instance-ids i-07e9c5779461b99d4 --query 'Reservations[*].Instances[*].PublicIpAddress' --output text
aws ec2 describe-instances --instance-ids [instanceId] --query 'Reservations[*].Instances[*].PublicIpAddress' --output text

# ssh login
# ssh -i workdayKeyPeir.pem ec2-user@18.176.48.143
ssh -i [/path/to/your-key.pem] ec2-user@[your-ec2-instance-public-ip]

# deploy
eb deploy
```


```bash

# eb-engine.logにアクセスする方法
# Elastic Beanstalk環境のEC2インスタンスにSSHで接続する
sudo less /var/log/eb-engine.log


```