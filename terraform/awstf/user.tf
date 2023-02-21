resource "aws_iam_user" "pxdeployuser" {
  name = format("%s-%s-user",var.name_prefix,var.config_name)
  path = "/"

  tags	= {
	Name 				= format("%s-%s",var.name_prefix,var.config_name)
	px-deploy_name 		= var.config_name
	px-deploy_username 	= var.PXDUSER
  }
  
}

resource "aws_iam_user_policy" "pxdeploypol" {
  name = format("%s-%s-policy",var.name_prefix,var.config_name)
  user = aws_iam_user.pxdeployuser.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
                "ec2:AttachVolume",
                "ec2:ModifyVolume",
                "ec2:DetachVolume",
                "ec2:CreateTags",
                "ec2:CreateVolume",
                "ec2:DeleteTags",
                "ec2:DeleteVolume",
                "ec2:DescribeTags",
                "ec2:DescribeVolumeAttribute",
                "ec2:DescribeVolumesModifications",
                "ec2:DescribeVolumeStatus",
                "ec2:DescribeVolumes",
                "ec2:DescribeInstances",
                "autoscaling:DescribeAutoScalingGroups",
                "s3:PutObject",
                "s3:GetObject",
                "s3:ListAllMyBuckets",
                "s3:CreateBucket",
                "s3:ListBucket",
                "s3:DeleteObject",
                "s3:GetBucketLocation"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}

resource "aws_iam_access_key" "pxdeploykey" {
  user = aws_iam_user.pxdeployuser.name
}