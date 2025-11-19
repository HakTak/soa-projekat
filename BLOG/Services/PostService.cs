using BLOG.Repositories;
using BLOG.Model;

namespace BLOG.Services
{
    public class PostService
    {
        private readonly IPostRepository _postRepository;

        public PostService(IPostRepository postRepository)
        {
            _postRepository = postRepository;
        }

        public Task<IEnumerable<Post>> GetPostsAsync()
        {
            return _postRepository.GetPostsAsync();
        }
        public Task<Post> GetPostByIdAsync(string id)
        {
            return _postRepository.GetPostByIdAsync(id);
        }

        public Task CreatePostAsync(Post Post)
        {
            return _postRepository.CreatePostAsync(Post);
        }

        public Task UpdatePostAsync(Post Post)
        {
            return _postRepository.UpdatePostAsync(Post);
        }

        public Task DeletePostAsync(string id)
        {
            return _postRepository.DeletePostAsync(id);
        }

        public async Task<Post?> TogglePostLikeAsync(string postId, string userId)
        {
            Post post = await _postRepository.GetPostByIdAsync(postId);
            if (post == null) { return null;}
            PostLike postLike = await _postRepository.GetPostLikeByUserIdAsync(userId, postId);
            if (postLike != null)
            {
                post.LikeCount --;
                await _postRepository.DeletePostLikeAsync(postLike.Id);
            }
            else
            {
                post.LikeCount ++;
                await _postRepository.CreatePostLikeAsync(new PostLike {PostId = post.Id, UserId = userId});
            }
            await _postRepository.UpdatePostAsync(post);
            return post;
        }

    }
}