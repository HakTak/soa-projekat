using BLOG.Models;
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

        public Task UpdateCommentAsync(Comment comment)
        {
            return _commentRepository.UpdateCommentAsync(comment);
        }

        public Task DeleteCommentAsync(string id)
        {
            return _commentRepository.DeleteCommentAsync(id);
        }
    }
}
