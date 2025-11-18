using BLOG.Controllers;
using BLOG.Model;

namespace BLOG.Repositories
{
    public interface IPostRepository
    {
        Task<IEnumerable<Post>> GetPostsAsync();
        Task<Post> GetPostByIdAsync(string id);

        Task<PostLike> GetPostLikeByUserIdAsync(string userId, string postId);
        Task CreatePostAsync(Post Post);
        Task CreatePostLikeAsync(PostLike postLike);
        Task DeletePostLikeAsync(string id);
        Task UpdatePostAsync(Post Post);
        Task DeletePostAsync(string id);
    }
}