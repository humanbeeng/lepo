<script lang="ts">
	let showMessage: boolean = false;
</script>

<body class="bg-black min-h-screen min-w-full">
	<nav class="navbar bg-black">
		<a class="btn btn-ghost text-white normal-case text-lg" href="/"
			>Lepo <span class="text-red-600">.</span></a
		>
	</nav>

	<!-- Hero -->
	<section class="grid grid-cols-1 px-2">
		<div class="pt-40 pb-20 flex flex-row justify-center">
			<h2 class=" text-4xl md:text-5xl font-black text-white text-center select-none">
				Lepo is an AI assistant which can answer <br />
				<span class="highlight">your queries</span>
				regarding
				<span class="highlight">your codebase</span>
			</h2>
		</div>
		<div class="flex flex-row justify-center">
			<h2 class="text-3xl font-bold text-zinc-400">Coming soon</h2>
		</div>

		<!-- form -->
		<!-- <div class="flex flex-row justify-center">
			{#if showMessage}
				<div class="text-center">
					<h3 class="font-extrabold text-3xl">Submitted</h3>
				</div>
			{:else}
				<div class="text-center">
					<form method="POST" class="">
						<div class="flex flex-cols-2 gap-x-4 justify-center">
							<input
								id="email"
								class="justify-center w-1/2 input lg:w-full bg-black form-control shadow-red-600 shadow-md backdrop-blur-lg text-zinc-400"
								aria-label="email"
								type="email"
								name="email"
								autocomplete="off"
								placeholder="Enter your email..."
								required
							/>
							<input
								type="submit"
								class="px-4 rounded-lg lg:px-5 bg-black shadow-red-600 form-control shadow-md backdrop-blur-lg text-zinc-400"
								value="Notify me!"
							/>
						</div>
					</form>
				</div>
			{/if}
		</div> -->
	</section>
</body>

<style>
	.highlight {
		background: url(//s2.svgbox.net/pen-brushes.svg?ic=brush-6&color=FF3232);
	}
</style>
