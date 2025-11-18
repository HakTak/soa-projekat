using System.Runtime.CompilerServices;
using BLOG.Model;
using MongoDB.Driver;
using MongoDatabaseSettings = BLOG.Database.MongoDatabaseSettings;
namespace BLOG.Repositories
{
    public class PostRepository : IPostRepository
    {
        private readonly IMongoCollection<Post> _posts;
        private readonly IMongoCollection<PostLike> _postLikes;

        public PostRepository(MongoDatabaseSettings settings)
        {
            var client = new MongoClient(settings.ConnectionString);
            var database = client.GetDatabase(settings.DatabaseName);
            _posts = database.GetCollection<Post>("posts");
            _postLikes = database.GetCollection<PostLike>("postLikes");
        }
        public async Task<IEnumerable<Post>> GetPostsAsync()
        {
            return await _posts
                .Find(_ => true)  
                .ToListAsync();
        }
        public async Task<Post> GetPostByIdAsync(string id)
        {
            return await _posts.Find(b => b.Id == id).FirstOrDefaultAsync();
        }
        public async Task CreatePostAsync(Post post)
        {
            await _posts.InsertOneAsync(post);
        }

        public async Task CreatePostLikeAsync(PostLike postLike)
        {
            await _postLikes.InsertOneAsync(postLike);
        }

        public async Task DeletePostLikeAsync(string id)
        {
            await _postLikes.DeleteOneAsync(b => b.Id == id);
        }
        public async Task DeletePostAsync(string id)
        {
            await _posts.DeleteOneAsync(b => b.Id == id);
        }

        public async Task UpdatePostAsync(Post post)
        {
            await _posts.ReplaceOneAsync(b => b.Id == post.Id, post);
        }

        public async Task<PostLike> GetPostLikeByUserIdAsync(string userId, string postId)
        {
            return await _postLikes.Find(pl => pl.UserId == userId & pl.PostId == postId).FirstOrDefaultAsync();
        }
    }
}