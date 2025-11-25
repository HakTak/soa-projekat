using BLOG.Model;
using BLOG.Services;
using Microsoft.AspNetCore.Mvc;

namespace BLOG.Controllers
{
    [ApiController]
    [Route("blog/[controller]")]
    public class PostController : ControllerBase
    {
        private readonly PostService _postSerivce;

        public PostController(PostService postService)
        {
            _postSerivce = postService;
        }

        [HttpGet("All")]
        public async Task<ActionResult<IEnumerable<Post>>> GetAll()
        {
            var posts = await _postSerivce.GetPostsAsync();
            if (!posts.Any())   
            {
                return NotFound("No posts found.");
            }
            return Ok(posts);
        }
        
        [HttpPost("toggleLike")]
        public async Task<ActionResult<Post>> ToggleLike([FromBody] LikeToggleRequest req)
        {
            var post = await _postSerivce.TogglePostLikeAsync(req.PostId, req.UserId);
            return Ok(post);
        } 

        [HttpPost("create")]
        public async Task<ActionResult<Post>> CreatePost([FromForm] PostDTO postDTO)
        {
            var post = new Post(postDTO);
            var uploadPath = Path.Combine(Directory.GetCurrentDirectory(), "wwwroot", "uploads"); 
            List<string> imageUrls = new List<string>();
            foreach (var file in postDTO.images)
            {
                var fileName = $"{Guid.NewGuid()}{Path.GetExtension(file.FileName)}";
                var filePath = Path.Combine(uploadPath, fileName);

                using (var stream = System.IO.File.Create(filePath))
                {
                    file.CopyTo(stream);
                }

                imageUrls.Add($"http:localhost:5251/wwwwroot/uploads/{fileName}"); // URL to serve
            }
            post.ImagePaths = imageUrls;
            await _postSerivce.CreatePostAsync(post);
            return Ok(post);
        }
    }
}
