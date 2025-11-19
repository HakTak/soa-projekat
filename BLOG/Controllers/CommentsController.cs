using BLOG.Model;
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

        [HttpGet("getAllByPostId/{postId}")] //svi komentari za dati post
        public async Task<ActionResult<IEnumerable<Comment>>> GetComments(string postId)
        {
            var comments = await _commentService.GetCommentsByPostIdAsync(postId);
            return Ok(comments);
        }

        [HttpGet("getSingleByCommentId/{id}")] //jedan komentar po id-u
        public async Task<ActionResult<Comment>> GetComment(string id)
        {
            var comment = await _commentService.GetCommentByIdAsync(id);
            if (comment == null) return NotFound();
            return Ok(comment);
        }

        [HttpPost("create")] //kreiranje komentara
        public async Task<ActionResult> CreateComment([FromBody] Comment comment)
        {
            await _commentService.CreateCommentAsync(comment);
            return CreatedAtAction(nameof(GetComment), new { id = comment.Id }, comment);
        }

        [HttpPut("edit/{id}")] //azuriranje komentara
        public async Task<ActionResult> UpdateComment(string id, [FromBody] Comment updatedComment)
        {
            await _commentService.UpdateCommentAsync(id, updatedComment);
            return NoContent();
        }

        [HttpDelete("delete/{id}")] //brisanje komentara
        public async Task<ActionResult> DeleteComment(string id)
        {
            await _commentService.DeleteCommentAsync(id);
            return NoContent();
        }
    }
}
