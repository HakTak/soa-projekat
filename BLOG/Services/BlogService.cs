using BLOG.Repositories;
using BLOG.Model;

namespace BLOG.Services
{
    public class BlogService
    {
        private readonly IBlogRepository _blogRepository;

        public BlogService(IBlogRepository BlogRepository)
        {
            _blogRepository = BlogRepository;
        }

        public Task<IEnumerable<Model.Blog>> GetBlogsAsync()
        {
            return _blogRepository.GetBlogAsync();
        }
        public Task<Model.Blog> GetPostByIdAsync(string id)
        {
            return _blogRepository.GetBlogByIdAsync(id);
        }

        public Task CreateBlogAsync(Model.Blog blog)
        {
            return _blogRepository.CreateBlogAsync(blog);
        }

        public Task UpdateBlogAsync(Model.Blog blog)
        {
            return _blogRepository.UpdateBlogAsync(blog);
        }

        public Task DeleteBlogAsync(string id)
        {
            return _blogRepository.DeleteBlogAsync(id);
        }

        public async Task<Model.Blog?> ToggleBlogLikeAsync(string BlogId, string userName)
        {
            Model.Blog blog = await _blogRepository.GetBlogByIdAsync(BlogId);
            if (blog == null) { return null;}
            BlogLike BlogLike = await _blogRepository.GetBlogLikeByUserIdAsync(userName, BlogId);
            if (BlogLike != null)
            {
                blog.LikeCount --;
                await _blogRepository.DeleteBlogLikeAsync(BlogLike.Id);
            }
            else
            {
                blog.LikeCount ++;
                await _blogRepository.CreateBlogLikeAsync(new BlogLike {BlogId = blog.Id, UserName = userName, LikedAt = new DateTime()});
            }
            await _blogRepository.UpdateBlogAsync(blog);
            return blog;
        }

    }
}