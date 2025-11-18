using BLOG.Model;
using BLOG.Repositories;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace BLOG.Services
{
    public class CommentService
    {
        private readonly ICommentRepository _commentRepository;

        public CommentService(ICommentRepository commentRepository)
        {
            _commentRepository = commentRepository;
        }

        public Task<IEnumerable<Comment>> GetCommentsByPostIdAsync(string postId)
        {
            return _commentRepository.GetCommentsByPostIdAsync(postId);
        }

        public Task<Comment> GetCommentByIdAsync(string id)
        {
            return _commentRepository.GetCommentByIdAsync(id);
        }

        public Task CreateCommentAsync(Comment comment)
        {
            return _commentRepository.CreateCommentAsync(comment);
        }

        public async Task<bool> UpdateCommentAsync(string id, Comment updatedComment)
        {
            var existingComment = await _commentRepository.GetCommentByIdAsync(id);
            if (existingComment == null)
            {
                return false; // Comment not found
            }

            existingComment.Text = updatedComment.Text;
            existingComment.UpdatedAt = System.DateTime.UtcNow;

            await _commentRepository.UpdateCommentAsync(existingComment);
            return true; // Update successful
            
        }

        public Task DeleteCommentAsync(string id)
        {
            return _commentRepository.DeleteCommentAsync(id);
        }
    }
}
