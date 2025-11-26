using System.Runtime.CompilerServices;
using BLOG.Model;
using MongoDB.Driver;
using MongoDatabaseSettings = BLOG.Database.MongoDatabaseSettings;
namespace BLOG.Repositories
{
    public class BlogRepository : IBlogRepository
    {
        private readonly IMongoCollection<Model.Blog> _blogs;
        private readonly IMongoCollection<BlogLike> _blogLikes;

        public BlogRepository(MongoDatabaseSettings settings)
        {
            var client = new MongoClient(settings.ConnectionString);
            var database = client.GetDatabase(settings.DatabaseName);
            _blogs = database.GetCollection<Model.Blog>("blogs");
            _blogLikes = database.GetCollection<BlogLike>("blogLikes");
        }
        public async Task<IEnumerable<Model.Blog>> GetBlogAsync()
        {
            return await _blogs
                .Find(_ => true)  
                .ToListAsync();
        }
        public async Task<Model.Blog> GetBlogByIdAsync(string id)
        {
            return await _blogs.Find(b => b.Id == id).FirstOrDefaultAsync();
        }
        public async Task<Model.Blog> CreateBlogAsync(Model.Blog blog)
        {
            await _blogs.InsertOneAsync(blog);
            return blog;
        }
        public async Task CreateBlogLikeAsync(BlogLike BlogLike)
        {
            await _blogLikes.InsertOneAsync(BlogLike);
        }

        public async Task DeleteBlogLikeAsync(string id)
        {
            await _blogLikes.DeleteOneAsync(b => b.Id == id);
        }
        public async Task DeleteBlogAsync(string id)
        {
            await _blogs.DeleteOneAsync(b => b.Id == id);
        }

        public async Task<Model.Blog> UpdateBlogAsync(Model.Blog post)
        {
            await _blogs.ReplaceOneAsync(b => b.Id == post.Id, post);
            return post;
        }

        public async Task<BlogLike> GetBlogLikeByUserIdAsync(string UserName, string blogId)
        {
            return await _blogLikes.Find(bl => bl.UserName == UserName & bl.BlogId == blogId).FirstOrDefaultAsync();
        }
    }
}