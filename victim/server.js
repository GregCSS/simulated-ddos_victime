// server.js — Simple Library Search backend
// Serves the static frontend (public/) and a small JSON API over data.json

const express = require('express');
const path = require('path');
const fs = require('fs');

const app = express();
const PORT = process.env.PORT || 3000;

// Load dataset once at startup
const dataPath = path.join(__dirname, 'data.json');
let books = [];
try {
  books = JSON.parse(fs.readFileSync(dataPath, 'utf-8'));
  console.log(`Loaded ${books.length} books from data.json`);
} catch (err) {
  console.error('Failed to load data.json:', err.message);
  process.exit(1);
}

// Serve the frontend
app.use(express.static(path.join(__dirname, 'public')));

// GET /api/books?q=orwell&genre=Sci-Fi&sort=rating&page=1&pageSize=16
app.get('/api/books', (req, res) => {
  const { q = '', genre = '', sort = 'title', page = '1', pageSize = '16' } = req.query;

  let results = books.filter((b) => {
    const query = q.trim().toLowerCase();
    const matchesQuery =
      !query ||
      b.title.toLowerCase().includes(query) ||
      b.author.toLowerCase().includes(query) ||
      b.genre.toLowerCase().includes(query);
    const matchesGenre = !genre || b.genre === genre;
    return matchesQuery && matchesGenre;
  });

  if (sort === 'title') results.sort((a, b) => a.title.localeCompare(b.title));
  if (sort === 'year') results.sort((a, b) => b.year - a.year);
  if (sort === 'rating') results.sort((a, b) => b.rating - a.rating);

  const pageNum = Math.max(1, parseInt(page, 10) || 1);
  const size = Math.max(1, parseInt(pageSize, 10) || 16);
  const total = results.length;
  const totalPages = Math.max(1, Math.ceil(total / size));
  const start = (pageNum - 1) * size;
  const pageItems = results.slice(start, start + size);

  res.json({
    total,
    page: pageNum,
    pageSize: size,
    totalPages,
    results: pageItems,
  });
});

// GET /api/genres — distinct genre list for the filter dropdown
app.get('/api/genres', (req, res) => {
  const genres = [...new Set(books.map((b) => b.genre))].sort();
  res.json(genres);
});

// GET /api/books/:id — single book detail
app.get('/api/books/:id', (req, res) => {
  const book = books.find((b) => b.id === parseInt(req.params.id, 10));
  if (!book) return res.status(404).json({ error: 'Book not found' });
  res.json(book);
});

app.listen(PORT, () => {
  console.log(`Library server running at http://localhost:${PORT}`);
});
