console.log("Menulis file...");
fs.writeFile("test.txt", "Halo dari Kilat!");
console.log("Membaca file:", fs.readFile("test.txt"));