const randomBytes = require('crypto').randomBytes;

const rideId = toUrlString(randomBytes(16));

function toUrlString(buffer) {
    return buffer.toString('base64')
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
}

console.log(rideId)