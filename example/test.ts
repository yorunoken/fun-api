async function main() {
   const cat = await fetch("http://localhost:3000/cat").then(res => res.json())
   console.log(cat.url);
}

main();
