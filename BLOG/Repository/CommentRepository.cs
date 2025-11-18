using BLOG.Models;
using MongoDB.Driver;
using BLOG.Database;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using MongoDatabaseSettings = BLOG.Database.MongoDatabaseSettings;

namespace BLOG.Repositories
{
    public class CommentRepository : ICommentRepository
    {
        private readonly IMongoCollection<Comment> _comments;

        public CommentRepository(MongoDatabaseSettings settings)
        {
            var client = new MongoClient(settings.ConnectionString);
            var database = client.GetDatabase(settings.DatabaseName);
            _comments = database.GetCollection<Comment>("comments");
        }

        public async Task<IEnumerable<Comment>> GetCommentsByPostIdAsync(string postId)
        {
            return await _comments.Find(c => c.PostId == postId).ToListAsync();
        }

        public async Task<Comment> GetCommentByIdAsync(string id)
        {
            return await _comments.Find(c => c.Id == id).FirstOrDefaultAsync();
        }

        public async Task CreateCommentAsync(Comment comment)
        {
            await _comments.InsertOneAsync(comment);
        }

        public async Task UpdateCommentAsync(Comment comment)
        {
            comment.UpdatedAt = DateTime.UtcNow;
            await _comments.ReplaceOneAsync(c => c.Id == comment.Id, comment);
        }

        public async Task DeleteCommentAsync(string id)
        {
            await _comments.DeleteOneAsync(c => c.Id == id);
        }
    }
}
