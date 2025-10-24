const userName = prompt("Enter your name:");

const canvas = document.getElementById('draw');
const ctx = canvas.getContext('2d');
canvas.width = 400;
canvas.height = 400;

let points = [];
let drawing = false;

canvas.addEventListener('mousedown', e => {
    drawing = true;
    points = [];
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.beginPath();
});

canvas.addEventListener('mousemove', e => {
    if (!drawing) return;
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    points.push({x, y});
    ctx.lineTo(x, y);
    ctx.strokeStyle = 'lime';
    ctx.lineWidth = 2;
    ctx.stroke();
});

canvas.addEventListener('mouseup', () => {
    drawing = false;
    ctx.beginPath();
    computeScore();
});

const leaderboardEl = document.getElementById('leaderboard');

function updateLeaderboard() {
    fetch('/leaderboard')
        .then(res => res.json())
        .then(data => {
            leaderboardEl.innerHTML = '';
            data.forEach(item => {
                const li = document.createElement('li');
                li.textContent = `${item.name}: ${item.score}`;
                leaderboardEl.appendChild(li);
            });
        });
}

function computeScore() {
    if (points.length < 10) return;
    const cx = points.reduce((a,p)=>a+p.x,0)/points.length;
    const cy = points.reduce((a,p)=>a+p.y,0)/points.length;
    const radii = points.map(p=>Math.hypot(p.x-cx,p.y-cy));
    const avg = radii.reduce((a,r)=>a+r,0)/radii.length;
    const variance = radii.reduce((a,r)=>a+(r-avg)**2,0)/radii.length;
    const normalized = Math.min(variance/avg**2,1);
    const score = Math.round((1 - normalized) * 100);
    document.getElementById('score').textContent = `Circle Accuracy: ${score}%`;

    fetch('/score', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({name: userName, score})
    }).then(updateLeaderboard);
}

window.onload = updateLeaderboard;
