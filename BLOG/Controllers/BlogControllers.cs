using Microsoft.AspNetCore.Mvc;

namespace BLOG.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class BlogController : ControllerBase
    {
        [HttpGet("health")]
        public IActionResult Health() => Ok(new { status = "Blog service is healthy" });

        [HttpGet]
        public IActionResult GetAll() => Ok(new[]
        {
            new { Id = 1, Title = "First Blog", Content = "Hello World" },
            new { Id = 2, Title = "Second Blog", Content = "Another Post" }
        });
    }
}

