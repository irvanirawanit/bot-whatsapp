
# Bot Whatsapp

Simple bot scan & send message 

Use at your own risk


## Requirment

 - [Golang 1.18](https://go.dev/doc/install)


## Deployment

To test this project run

```bash
  go run .
```


## API Reference

#### Show QR Code on webpage (additionally QR Code show in terminal)

```http
  GET /
```

#### Send Message

```http
  POST /message
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `phone` | `string` | **Required**. example "628123456789" |
| `message` | `string` | **Required**. example "Your OTP Code 123456" |


## Support

For support, email irvanirawanit@gmail.com


## Credit

package [tulir/whatsmeow](https://github.com/tulir/whatsmeow)

