console.log("Menulis file...");
fs.writeFile("test.txt", "Halo dari Jolt!");
console.log("Membaca file:", fs.readFile("test.txt"));