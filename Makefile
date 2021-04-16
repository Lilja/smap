text-proxyjump:
	go build -o main && ./main --verbose --columns host,port,username,jump,nojumps --file test/proxy
