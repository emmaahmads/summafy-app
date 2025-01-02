docker build -t summafy-repo .
docker tag summafy-repo:latest 120569638976.dkr.ecr.us-east-1.amazonaws.com/summafy-repo:latest
docker push 120569638976.dkr.ecr.us-east-1.amazonaws.com/summafy-repo:latest
