psql postgresql://emma:happybirthday@my-rds.cnieu4oowqof.us-east-1.rds.amazonaws.com:5432/summafy?sslmode=allow

ssh -N -L 5436:summafy.my-rds.cnieu4oowqof.us-east-1.rds.amazonaws.com:5432 -p 22 -i mykey.pem ec2-user@107.23.208.53