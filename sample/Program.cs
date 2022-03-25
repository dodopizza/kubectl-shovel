using System.IO;
using Microsoft.AspNetCore.Hosting;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;

namespace Sample
{
	internal static class Program
	{
		private static void Main(string[] args) => CreateHostBuilder(args)
			.ConfigureLogging(logging => logging
				.AddConsole()
				.SetMinimumLevel(LogLevel.Debug))
			.ConfigureWebHost(host => host
				.UseKestrel((ctx, options) => options.Configure(ctx.Configuration.GetSection("Kestrel")))
				.UseStartup<Startup>())
			.Build()
			.Run();

		private static IHostBuilder CreateHostBuilder(string[] args) => new HostBuilder()
			.ConfigureAppConfiguration(conf => conf
				.AddCommandLine(args)
				.SetBasePath(Directory.GetCurrentDirectory())
				.AddJsonFile("appsettings.json", false, false)
				.AddEnvironmentVariables());
	}
}