const result = net.fetch("https://api.shny.my.id/api/status");
console.log("Status:", result.status);
console.log("Body:", result.body);