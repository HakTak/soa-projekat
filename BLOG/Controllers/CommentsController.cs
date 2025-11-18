using BLOG.Models;
using BLOG.Services;
using Microsoft.AspNetCore.Mvc;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace BLOG.Controllers
{
    [Route("blog/[controller]")]
    [ApiController]
    public class CommentsController : ControllerBase
    {
        private readonly CommentService _commentService;

        public CommentsController(CommentService commentService)
        {
            _commentService = commentService;
        }

        [HttpGet("{postId}")]
        public async Task<ActionResult<IEnumerable<Comment>>> GetComments(string postId)
        {
            var comments = await _commentService.GetCommentsByPostIdAsync(postId);
            return Ok(comments);
        }

        [HttpGet("single/{id}")]
        public async Task<ActionResult<Comment>> GetComment(string id)
        {
            var comment = await _commentService.GetCommentByIdAsync(id);
            if (comment == null) return NotFound();
            return Ok(comment);
        }

        [HttpPost]
        public async Task<ActionResult> CreateComment([FromBody] Comment comment)
        {
            await _commentService.CreateCommentAsync(comment);
            return CreatedAtAction(nameof(GetComment), new { id = comment.Id }, comment);
        }

        [HttpPut("{id}")]
        public async Task<ActionResult> UpdateComment(string id, [FromBody] Comment updatedComment)
        {
            var existingComment = await _commentService.GetCommentByIdAsync(id);
            if (existingComment == null) return NotFound();

            existingComment.Text = updatedComment.Text;
            existingComment.UpdatedAt = System.DateTime.UtcNow;

            await _commentService.UpdateCommentAsync(existingComment);
            return NoContent();
        }

        [HttpDelete("{id}")]
        public async Task<ActionResult> DeleteComment(string id)
        {
            await _commentService.DeleteCommentAsync(id);
            return NoContent();
        }
    }
}
