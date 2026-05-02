<?php
require_once __DIR__ . '/cache.php';
set_cache_headers(__FILE__);
?>

<!DOCTYPE html>
<html lang="de">

<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Masseverbrauch Rechner</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}

		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
			background: #f5f5f5;
			padding: 1rem;
			min-height: 100vh;
		}

		.container {
			max-width: 600px;
			margin: 0 auto;
			background: white;
			padding: 1.5rem;
			border-radius: 8px;
			box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
		}

		h1 {
			font-size: 1.5rem;
			margin-bottom: 1.5rem;
			color: #333;
		}

		h2 {
			font-size: 1.1rem;
			margin: 1rem 0 0.5rem;
			color: #555;
		}

		.form-group {
			margin-bottom: 1rem;
		}

		label {
			display: block;
			margin-bottom: 0.25rem;
			color: #666;
			font-size: 0.9rem;
		}

		input[type="number"] {
			width: 100%;
			padding: 0.5rem;
			border: 1px solid #ddd;
			border-radius: 4px;
			font-size: 1rem;
		}

		.slot-grid {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 0.5rem;
			margin-bottom: 0.5rem;
		}

		.slot-grid label {
			font-size: 0.8rem;
		}

		.btn-group {
			display: flex;
			gap: 0.5rem;
			margin-top: 1.5rem;
		}

		button {
			flex: 1;
			padding: 0.75rem;
			border: none;
			border-radius: 4px;
			font-size: 1rem;
			cursor: pointer;
		}

		button[type="submit"] {
			background: #007bff;
			color: white;
		}

		button[type="reset"] {
			background: #6c757d;
			color: white;
		}

		.result-row {
			display: flex;
			justify-content: space-between;
			padding: 0.5rem 0;
			border-bottom: 1px solid #eee;
		}

		.result-row strong {
			color: #333;
		}

		dialog {
			border: none;
			border-radius: 8px;
			padding: 1.5rem;
			max-width: 500px;
			width: 90%;
			position: fixed;
			top: 50%;
			left: 50%;
			transform: translate(-50%, -50%);
			margin: 0;
		}

		dialog::backdrop {
			background: rgba(0, 0, 0, 0.5);
		}

		.close-btn {
			background: #dc3545;
			color: white;
			margin-top: 1rem;
			width: 100%;
		}

		@media (max-width: 480px) {
			body {
				padding: 0.5rem;
			}

			.container {
				padding: 1rem;
			}

			.slot-grid {
				grid-template-columns: 1fr;
			}
		}
	</style>
</head>

<body>
	<div class="container">
		<h1>Masseverbrauch Rechner</h1>
		<div style="text-align: right; font-size: 0.8rem; color: #999;">
			Web: <?php echo rtrim(file_get_contents(__DIR__ . '/VERSION') ?: 'dev'); ?>
			(API: <span id="api-version">...</span>)
		</div>

		<form hx-post="calculate.php"
			hx-target="#results"
			hx-swap="innerHTML">
			<h2>Hauptmasse</h2>

			<div class="form-group">
				<label for="main-flow">Flow (kg/s)</label>
				<input type="number" id="main-flow" name="flow" step="0.1" min="0" required>
			</div>

			<div class="form-group">
				<label for="main-runtime">Laufzeit (Minuten)</label>
				<input type="number" id="main-runtime" name="runtime_minutes" step="0.01" min="0" required>
			</div>

			<h2>Flakes - Turmpositionen</h2>

			<?php for ($i = 1; $i <= 5; $i++): ?>
				<div class="form-group">
					<label>Turmposition <?php echo $i; ?></label>
					<div class="slot-grid">
						<div>
							<label for="tower-slot-<?php echo $i; ?>-before">KG vor</label>
							<input type="number" id="tower-slot-<?php echo $i; ?>-before"
								name="slots[<?php echo $i - 1; ?>][before]" step="0.01" min="0">
						</div>
						<div>
							<label for="tower-slot-<?php echo $i; ?>-after">KG nach</label>
							<input type="number" id="tower-slot-<?php echo $i; ?>-after"
								name="slots[<?php echo $i - 1; ?>][after]" step="0.01" min="0">
						</div>
					</div>
				</div>
			<?php endfor; ?>

			<div class="btn-group">
				<button type="reset">Zurücksetzen</button>
				<button type="submit">Berechnen</button>
			</div>
		</form>

		<div id="results"></div>
	</div>

	<script>
		(function() {
			const apiBase = ('<?php echo getenv("API_URL") ?: "http://localhost:8080/api/calculate"; ?>').replace('/api/calculate', '');
			fetch(apiBase + '/api/version')
				.then(r => r.json())
				.then(d => document.getElementById('api-version').textContent = d.version)
				.catch(() => document.getElementById('api-version').textContent = 'n/a');
		})();
	</script>
</body>

</html>
