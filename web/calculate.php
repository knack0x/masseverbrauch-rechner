<?php
# Cache Control
require_once __DIR__ . "/cache.php";
set_cache_headers(__FILE__);
?>

<?php
$apiUrl = 'http://localhost:8080/api/calculate';

// Read POST data
$flow = floatval($_POST['flow'] ?? 0);
$runtimeMinutes = floatval($_POST['runtime_minutes'] ?? 0);

$slots = [];
if (isset($_POST['slots'])) {
	foreach ($_POST['slots'] as $slot) {
		$slots[] = [
			'before' => floatval($slot['before'] ?? 0),
			'after' => floatval($slot['after'] ?? 0)
		];
	}
}

// Prepare API request
$data = [
	'flow' => $flow,
	'runtime_minutes' => $runtimeMinutes,
	'slots' => $slots
];

$ch = curl_init($apiUrl);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);
curl_setopt($ch, CURLOPT_POSTFIELDS, json_encode($data));

$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

if ($httpCode !== 200) {
	echo '<div style="color: red; padding: 1rem;">Fehler bei der Berechnung</div>';
	exit;
}

$result = json_decode($response, true);
?>

<dialog id="results-dialog">
	<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
		<h2 style="margin: 0;">Verbrauch</h2>
	</div>

	<div class="result-row">
		<strong>Hauptmasse</strong>
	</div>
	<div class="result-row">
		<span><?php echo number_format($result['hauptmasse_kg'], 2, ',', '.'); ?> kg</span>
		<span><?php echo number_format($result['hauptmasse_percent'], 2, ',', '.'); ?>%</span>
	</div>

	<div class="result-row" style="margin-top: 1rem;">
		<strong>Zusatzmasse</strong>
	</div>

	<?php
	function trSlotName(string $in): string
	{
		$tr = [
			"Tower Slot 1" => "Turmposition 1",
			"Tower Slot 2" => "Turmposition 2",
			"Tower Slot 3" => "Turmposition 3",
			"Tower Slot 4" => "Turmposition 4",
			"Tower Slot 5" => "Turmposition 5",
		];
		return $tr[$in] ?? $in;
	}

	foreach ($result['slots'] as $slot):
	?>
		<div class="result-row">
			<span><?php echo htmlspecialchars(trSlotName($slot['name'])); ?></span>
			<span><?php echo number_format($slot['kg'], 2, ',', '.'); ?> kg (<?php echo number_format($slot['percent'], 2, ',', '.'); ?>%)</span>
		</div>
	<?php endforeach; ?>

	<div class="result-row" style="margin-top: 1rem; border-top: 2px solid #333; padding-top: 0.5rem;">
		<strong>Gesamt: <?php echo number_format($result['total_kg'], 2, ',', '.'); ?> kg</strong>
	</div>

	<button class="close-btn" onclick="this.closest('dialog').close()">Schließen</button>
</dialog>

<script>
	document.getElementById('results-dialog').showModal();
</script>
