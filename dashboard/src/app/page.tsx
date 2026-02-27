export default function OverviewPage() {
	return (
		<div>
			<h2 className="mb-6 text-2xl font-bold text-white">Cluster Overview</h2>

			<div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
				<StatCard label="Nodes" value="--" />
				<StatCard label="Containers" value="--" />
				<StatCard label="Deployments" value="--" />
				<StatCard label="Services" value="--" />
			</div>

			<p className="mt-8 text-sm text-gray-500">
				Dashboard will be populated once the API is connected.
			</p>
		</div>
	);
}

function StatCard({ label, value }: { label: string; value: string }) {
	return (
		<div className="rounded-lg border border-gray-800 bg-gray-900 p-6">
			<p className="text-sm font-medium text-gray-400">{label}</p>
			<p className="mt-2 text-3xl font-bold text-white">{value}</p>
		</div>
	);
}
