TF_FILES=$(wildcard *.tf)
VAR_FILES=$(wildcard *.tfvars)
PY_FILES=$(wildcard *.py)

default: terraform.plan

terraform.plan: $(TF_FILES) $(VAR_FILES) $(PY_FILES)
	terraform plan -out terraform.plan

plan: terraform.plan

apply: terraform.plan
	terraform apply terraform.plan || rm terraform.plan

clean:
	rm -rf terraform.plan



import:
	terraform import aws_api_gateway_integration.gateway_to_sqs zj2i6wj4d5/d7t6bb/POST