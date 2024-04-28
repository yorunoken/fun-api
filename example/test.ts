async function main() {
   const player = await fetch("http://localhost:3000/api/user?username=yorunoken").then(res => res.json())
   console.log(player.statistics);
}

main();
