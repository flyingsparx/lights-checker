# Hue Bulb Checker

Automatically check-up on bulb status at night and alert if unusual behaviour.

Potentially useful for knowing when lights were turned on when they're not supposed to be.

The program uses [Mailgun](http://mailgun.com) to send email notifications. [GoCron](http://github.com/jasonlvhit/gocron) is used as an internal scheduler.

## Features

* Checks status of all bulbs connected to the specified Hue bridge every minute between specified hours.
* If there is are bulbs turned on, then send an email notification to the specified address.
* If bulbs are turned on for consecutive minutes, multiple emails will not be sent.

## Preparing

* Obtain a [Hue bridge username](https://developers.meethue.com/documentation/configuration-api#71_create_user) for your bridge
* Setup a Mailgun account, configure your domain for sending through Mailgun, and obtain your private and public key.
* Create a `config.json` with the required information using the skeleton file supplied.

## Using the program

Run `go get` to download the Mailgun and GoCron dependencies, and then `go run light_checker`.
