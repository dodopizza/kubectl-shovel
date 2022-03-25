using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Diagnostics.HealthChecks;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Diagnostics.HealthChecks;
using Microsoft.Extensions.Hosting;

namespace Sample
{
	internal sealed class Startup
	{
		private IConfiguration _configuration;

		public Startup(IConfiguration configuration)
		{
			_configuration = configuration;
		}

		public void ConfigureServices(IServiceCollection services)
		{
			services.AddRouting();
			services.AddMvcCore();
			services
				.AddHealthChecks()
				.AddCheck("live", () => HealthCheckResult.Healthy());
		}

		public void Configure(IApplicationBuilder app, IHostEnvironment environment)
		{
			app.UseDeveloperExceptionPage();
			app.UseRouting();
			app.UseEndpoints(endpoints =>
				{
					endpoints.MapControllers();
					endpoints.MapHealthChecks("/health/live", new HealthCheckOptions());
				}
			);
		}
	}
}