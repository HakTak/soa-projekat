using BLOG.Models;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace BLOG.Repositories
{
    public interface ICommentRepository
    {
        Task<IEnumerable<Comment>> GetCommentsByPostIdAsync(string postId);
        Task<Comment> GetCommentByIdAsync(string id);
        Task CreateCommentAsync(Comment comment);
        Task UpdateCommentAsync(Comment comment);
        Task DeleteCommentAsync(string id);
    }
}
