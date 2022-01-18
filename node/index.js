const https = require('https')

function doAThing() {
  throw new Error("KABOOM")
}

let apiTokenID = ''
let apiTokenSecret = ''

try {
  doAThing()
} catch (e) {
  // This could probably be simpler, but I haven't written node in years and it was only a tiny thing at that
  // Using stdlib always seems like a good idea to me. Everyone knows stdlib, right?
  const data = JSON.stringify({
    timestamp: new Date(),
    text: e.toString(),
    tags: ["tinyhatchet", "node", "example"]
  })

  const options = {
    hostname: 'tinyhatchet.com',
    port: 443,
    path: '/ingest.json',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Content-Length': data.length,
      'Authorization': 'Basic ' + Buffer.from(apiTokenID + ':' + apiTokenSecret).toString('base64')
    }
  }

  const req = https.request(options, res => {
    console.log(`statusCode: ${res.statusCode}`)

    res.on('data', d => {
      process.stdout.write(d)
    })
  })

  req.on('error', error => {
    console.error(error)
  })

  req.write(data)
  req.end()
}