
using BLOG.Model;

namespace BLOG.Repositories
{
    public interface IBlogRepository
    {
        Task<IEnumerable<Model.Blog>> GetBlogAsync();
        Task<Model.Blog> GetBlogByIdAsync(string id);

        Task<BlogLike> GetBlogLikeByUserIdAsync(string userId, string postId);
        Task<Model.Blog> CreateBlogAsync(Model.Blog post);
        Task CreateBlogLikeAsync(BlogLike blogLike);
        Task DeleteBlogLikeAsync(string id);
        Task<Model.Blog> UpdateBlogAsync(Model.Blog blog);
        Task DeleteBlogAsync(string id);
    }
}