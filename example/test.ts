async function main() {
   const cat = await fetch("http://localhost:3000/api/token").then(res => res.text())
   console.log(cat);
}

main();
