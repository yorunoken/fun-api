window.onload = async function () {
    const queryParams = new URLSearchParams(window.location.search);

    const username = queryParams.get("username") ?? "yorunoken";
    const mode = queryParams.get("mode") ?? "osu";

    const player = await fetch(`/api/user?username=${username}&mode=${mode}`).then((res) => res.json());

    const { statistics } = player;

    console.log(player);

    if (!player.id) {
        document.title = `User not found!`;
        console.error("either user was not found, or the request was bad.");
        return;
    }

    document.title = `User card for ${player.username} (#${statistics.global_rank.toLocaleString()})`;
    const metaDescription = document.querySelector('meta[name="description"]');
    if (metaDescription) {
        metaDescription.setAttribute("content", `User card generated by yorunoken, for ${player.username}`);
    } else {
        const newMeta = document.createElement("meta");
        newMeta.setAttribute("name", "description");
        newMeta.setAttribute("content", `User card generated by yorunoken, for ${player.username}`);
        document.head.appendChild(newMeta);
    }
    console.log("updated shiii");

    document.getElementById("grade-ssh").textContent = statistics.grade_counts.ssh;
    document.getElementById("grade-ss").textContent = statistics.grade_counts.ss;
    document.getElementById("grade-sh").textContent = statistics.grade_counts.sh;
    document.getElementById("grade-s").textContent = statistics.grade_counts.s;
    document.getElementById("grade-a").textContent = statistics.grade_counts.a;

    document.getElementById("username").textContent = player.username;
    document.getElementById("rank").textContent = `#${statistics.global_rank.toLocaleString()}`;
    document.getElementById("country-rank").textContent = `#${statistics.country_rank.toLocaleString()}`;
    document.getElementById("pp").textContent = `${statistics.pp.toFixed(2)}`;
    document.getElementById("accuracy").textContent = `${statistics.hit_accuracy.toFixed(2)}%`;
    document.getElementById("score").textContent = `${statistics.ranked_score.toLocaleString()}`;
    document.getElementById("playcount").textContent = `${statistics.play_count.toLocaleString()}`;
    document.getElementById("combo").textContent = `${statistics.maximum_combo.toLocaleString()}`;
    document.getElementById("avatar").src = player.avatar_url;

    document.getElementById("graph").src = `/api/graph?points=${player.rankHistory.data.join(",")}`;

    const { level } = statistics;
    document.getElementById("level").textContent = `${level.current}.${level.progress.toFixed()}%`;

    console.log(level.progress.toFixed());
    document.getElementById("level-bar").style.background = `linear-gradient(to right, #5C99AB ${(level.progress + 0.5).toFixed()}%, #2F393E ${(level.progress + 0.5).toFixed()}%)`;
};
