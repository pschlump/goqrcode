
all:
	go build

test: test1 test3 test4
	@echo PASS

test1:
	qrcode -o out/test1.png "otpauth://totp/www.2c-why.com:pschlump@gmail.com?secret=4S62BZNFXXSZLCRO&issuer=www.2c-why.com"
	diff out/test1.png ref/test1.png

test2:
	qrcode -o test2.png \
		"{\"address\":\"1d217eaaaac1deb2e75d1ec44bcae03a1227a126\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"6a0c48361aaa0c8bbcb33d7b53bef982b6620c7b1e5fd1d1c24457fc4416f517\",\"cipherparams\":{\"iv\":\"305d1eb07d717eaaa3668faeb7d04c43\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"8d27d8bb2ec646bfd577f833b72a461965769faef7fd5daf70b0c80857ffc589\"},\"mac\":\"6bbfb5cab3aed19070b7927fccfc62a56452fdc2a1325f70df23ea8c51794382\"},\"id\":\"e9c6ccb4-b1e2-45e5-bfca-7d39004cb3f4\",\"version\":3}"

test3:
	qrcode -txt -o out/test3.txt "otpauth://totp/www.2c-why.com:pschlump@gmail.com?secret=4S62BZNFXXSZLCRO&issuer=www.2c-why.com"
	diff out/test3.txt ref/test3.txt

test4:
	qrcode -svg -o out/test4.svg "otpauth://totp/www.2c-why.com:pschlump@gmail.com?secret=4S62BZNFXXSZLCRO&issuer=www.2c-why.com"
	diff out/test4.svg ref/test4.svg

