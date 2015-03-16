<?php

define('SENDER_NAME', 'BOSBEC');
define('BASE_URL', 'https://api.mobileresponse.se/');

$username = $_ENV['MOBILERESPONSE_API_USERNAME'];
$password = $_ENV['MOBILERESPONSE_API_PASSWORD'];
$message = 'Your message';
$recipients = array('Your phone number');

if (empty($username)) {
    die('Missing MobileResponse API password. Did you forget to set the MOBILERESPONSE_API_PASSWORD environment variable?' . PHP_EOL);
}

if (empty($password)) {
    die('Missing MobileResponse API password. Did you forget to set the MOBILERESPONSE_API_PASSWORD environment variable?' . PHP_EOL);
}

if (!function_exists('curl_init')) {
    die('Missing function curl_init. Is cURL installed as a PHP extension?' . PHP_EOL);
}

function authenticate($username, $password) {
    $request = array(
        'data' => array(
            'username' => $username,
            'password' => $password
        )
    );

    $curl = curl_init(BASE_URL . 'authenticate');

    curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($curl, CURLOPT_POST, true);
    curl_setopt($curl, CURLOPT_POSTFIELDS, json_encode($request));

    $content = curl_exec($curl);

    if (curl_errno($curl) > 0) {
        die(curl_error($curl) . PHP_EOL);
    }

    curl_close($curl);

    $response = json_decode($content);

    $authentication_token = $response->data->id;

    return $authentication_token;
}

function send_sms($username, $password, $message, $recipients) {
    if (!is_array($recipients)) {
        $recipients = array($recipients);
    }

    $request = array(
        'data' => array(
            'username' => $username,
            'password' => $password,
            'recipients' => $recipients,
            'message' => $message,
            'senderName' => SENDER_NAME
        )
    );

    $curl = curl_init(BASE_URL . 'quickie/send-message');

    curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($curl, CURLOPT_POST, true);
    curl_setopt($curl, CURLOPT_POSTFIELDS, json_encode($request));

    curl_exec($curl);

    if (curl_errno($curl) > 0) {
        die(curl_error($curl) . PHP_EOL);
    }

    curl_close($curl);
}

function is_authenticated($authentication_token) {
    $request = array(
        'data' => array(),
        'authenticationToken' =>  $authentication_token
    );

    $curl = curl_init(BASE_URL . 'is-authenticated');

    curl_setopt($curl, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($curl, CURLOPT_POST, true);
    curl_setopt($curl, CURLOPT_POSTFIELDS, json_encode($request));

    $content = curl_exec($curl);

    if (curl_errno($curl) > 0) {
        die(curl_error($curl) . PHP_EOL);
    }

    curl_close($curl);

    $response = json_decode($content);

    return $response->status == 'Success';
}

echo 'Sending "' . $message . '" to ' . implode(', ', $recipients) . PHP_EOL;

send_sms($username, $password, $message, $recipients);

echo 'Authenticating' . PHP_EOL;

$authentication_token = authenticate($username, $password);

echo 'Authenticated and received this "' . $authentication_token . '" authentication token' . PHP_EOL;

echo 'Checking if the authentication token is still valid' . PHP_EOL;

$is_authenticated = is_authenticated($authentication_token);

if ($is_authenticated) {
    echo '"' . $authentication_token . '" is valid' . PHP_EOL;
} else {
    echo '"' . $authentication_token . '" is not valid' . PHP_EOL;
}