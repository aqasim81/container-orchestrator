import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
	title: "Container Orchestrator",
	description: "Kubernetes-style container orchestrator management dashboard",
};

const navItems = [
	{ href: "/", label: "Overview" },
	{ href: "/nodes", label: "Nodes" },
	{ href: "/containers", label: "Containers" },
	{ href: "/deployments", label: "Deployments" },
	{ href: "/services", label: "Services" },
];

export default function RootLayout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en">
			<body className="bg-gray-950 text-gray-100 antialiased">
				<div className="flex min-h-screen">
					{/* Sidebar */}
					<aside className="w-64 border-r border-gray-800 bg-gray-900 p-6">
						<h1 className="mb-8 text-lg font-bold tracking-tight text-white">Orchestrator</h1>
						<nav className="space-y-1">
							{navItems.map((item) => (
								<a
									key={item.href}
									href={item.href}
									className="block rounded-md px-3 py-2 text-sm font-medium text-gray-400 transition-colors hover:bg-gray-800 hover:text-white"
								>
									{item.label}
								</a>
							))}
						</nav>
					</aside>

					{/* Main content */}
					<main className="flex-1 p-8">{children}</main>
				</div>
			</body>
		</html>
	);
}
